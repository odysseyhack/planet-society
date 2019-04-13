//
//  TransactionOverViewViewController.swift
//  Permission Hub
//
//  Created by Corné on 13/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionFlowViewController: UIPageViewController {

    // MARK: - Private properties

    private let steps = TransactionFlowStep.allCases
    private var currentStepIndex = 0

    private let service = NetworkingService()

    // MARK: - Initialization

    init(steps: [TransactionFlowStep]) {
        self.steps = steps

        super.init(
            transitionStyle: .scroll,
            navigationOrientation: .horizontal,
            options: [:])
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        view.backgroundColor = PHColors.lightGray

        // configure navigation bar
        navigationController?.navigationBar.tintColor = PHColors.greyishBrown
        navigationItem.backBarButtonItem = UIBarButtonItem(title: "", style: .plain, target: nil, action: nil)

        // set initial viewcontroller
        let firstViewControllerOrNil = self.steps.map { $0.viewController }.first
        if let viewController = firstViewControllerOrNil {

            setViewControllers(
                [viewController],
                direction: .forward,
                animated: false)
        }
    }
}

final class TransactionOverViewViewController: PHTableViewController {

    // MARK: - Private properties

    private lazy var bottomStackView: UIStackView = {

        let stackView = UIStackView()
        stackView.translatesAutoresizingMaskIntoConstraints = false

        stackView.axis = .horizontal
        stackView.alignment = .center
        stackView.distribution = .equalCentering
        stackView.spacing = 20

        return stackView
    }()

    private lazy var declineButton: UIButton = {

        let button = UIButton()
        button.translatesAutoresizingMaskIntoConstraints = false

        button.widthAnchor.constraint(equalToConstant: 92).isActive = true
        button.heightAnchor.constraint(equalToConstant: 34).isActive = true

        button.titleLabel?.font = PHFonts.regular()

        button.setTitleColor(PHColors.greyishBrown, for: .normal)
        button.backgroundColor = .white

        button.layer.cornerRadius = 5
        button.layer.borderWidth = 1
        button.layer.borderColor = PHColors.greyishBrown.cgColor

        button.setTitle("Decline", for: .normal)

        button.addTarget(
            self,
            action: #selector(declineButtonTapped),
            for: .touchUpInside)

        return button
    }()

    private lazy var continueButton: UIButton = {

        let button = UIButton()
        button.translatesAutoresizingMaskIntoConstraints = false

        button.widthAnchor.constraint(equalToConstant: 92).isActive = true
        button.heightAnchor.constraint(equalToConstant: 34).isActive = true

        button.titleLabel?.font = PHFonts.regular()

        button.setTitleColor(PHColors.topaz.withAlphaComponent(0.5), for: .disabled)
        button.setTitleColor(.white, for: .normal)
        button.backgroundColor = .white

        button.layer.cornerRadius = 5
        button.layer.borderWidth = 1
        button.layer.borderColor = PHColors.topaz.cgColor

        button.setTitle("Continue", for: .normal)

        button.addTarget(
            self,
            action: #selector(continueButtonTapped),
            for: .touchUpInside)

        // disabled by default
        button.isEnabled = false

        return button
    }()

    private var transaction: TransactionNotification

    // MARK: - Life cycle

    init(transaction: TransactionNotification) {
        self.transaction = transaction

        var items = [PHTableViewViewCellType]()
        items.append(.notification(
            type: .verification,
            text: "This company is verified"))
        items.append(.notification(
            type: .warning,
            text: "Permission warning!"))
        items.append(.description(
            date: transaction.date,
            title: transaction.title,
            description: transaction.description))

        transaction.items.forEach {
            items.append(.transactionItem(item: $0))
        }

        super.init(
            title: transaction.requesterName,
            items: items)
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    override func viewDidLoad() {
        super.viewDidLoad()

        // set selection delegate to self
        delegate = self

        view.addSubview(bottomStackView)
        bottomStackView.heightAnchor.constraint(equalToConstant: 75).isActive = true
        bottomStackView.topAnchor.constraint(equalTo: super.tableView.bottomAnchor).isActive = true
        bottomStackView.centerXAnchor.constraint(equalTo: view.centerXAnchor).isActive = true
        bottomStackView.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor).isActive = true

        bottomStackView.addArrangedSubview(declineButton)
        bottomStackView.addArrangedSubview(continueButton)
    }

    // MARK: - Networking

    private func respondToTransaction(isAccepted: Bool) {

        do {
            let service = NetworkingService.shared
            try service.respondToTransaction(
                withId: transaction.transactionID,
                isAccepted: true) { [unowned self] response in

                    self.dismiss(animated: true)

                    switch response {
                    case .failure(let error):
                        print(error)
                    default: break
                    }
            }
        } catch {
            print(error)
        }
    }

    // MARK: - Selectors

    @objc private func declineButtonTapped(_ sender: UIButton) {
        respondToTransaction(isAccepted: false)
    }

    @objc private func continueButtonTapped(_ sender: UIButton) {
        respondToTransaction(isAccepted: true)
    }

    // MARK: - Helpers

    private func validateSelection() {

        let areAllItemsAccepted = transaction.items.filter { !$0.isAccepted }.count == 0
        continueButton.isEnabled = areAllItemsAccepted
    }
}

extension TransactionOverViewViewController: PHTableViewControllerDelegate {

    func didSelect(item: PHTableViewViewCellType) {

        switch item {
        case .notification(let type, let text):

            let items: [PHTableViewViewCellType]
            switch type {
            case .warning:
                items = transaction.analysis.map { .description(
                    date: Date(),
                    title: "",
                    description: $0) }
                let viewController = PHTableViewController(
                    title: text,
                    items: items)
                navigationController?.pushViewController(viewController, animated: true)

            case .verification:
                 items = transaction.verification.map { .description(
                    date: Date(),
                    title: "",
                    description: $0) }
            }

            let viewController = PHTableViewController(
                title: text,
                items: items)
            navigationController?.pushViewController(viewController, animated: true)

        case .transactionItem(let item):
            if let index = self.transaction.items.index(of: item) {
                transaction.items[index].isAccepted = true
            }
            validateSelection()

        default:
            break
        }
    }

    func didDeselect(item: PHTableViewViewCellType) {

        switch item {
        case .transactionItem(let item):
            if let index = self.transaction.items.index(of: item) {
                transaction.items[index].isAccepted = false
            }
            validateSelection()

        default:
            break
        }
    }
}
