//
//  TransactionFlowViewController.swift
//  Permission Hub
//
//  Created by Corné on 13/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

final class TransactionFlowViewController: UIViewController {

    // MARK: - Private properties

    private lazy var pageViewController: UIPageViewController = {

        let pageViewController = UIPageViewController(
            transitionStyle: .scroll,
            navigationOrientation: .horizontal)

        return pageViewController
    }()

    private lazy var bottomStackView: UIStackView = {

        let stackView = UIStackView()
        stackView.translatesAutoresizingMaskIntoConstraints = false

        stackView.axis = .horizontal
        stackView.alignment = .center
        stackView.distribution = .equalCentering
        stackView.spacing = 20

        return stackView
    }()

    private lazy var optionsButton: UIButton = {

        let button = UIButton()
        button.translatesAutoresizingMaskIntoConstraints = false

        let dimension: CGFloat = 44
        button.widthAnchor.constraint(equalToConstant: dimension).isActive = true
        button.heightAnchor.constraint(equalToConstant: dimension).isActive = true

        let image = UIImage(named: "question_mark")
        button.setImage(image, for: .normal)

        button.addTarget(
            self,
            action: #selector(optionsButtonTapped),
            for: .touchUpInside)

        return button
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

        button.setTitle("Cancel", for: .normal)

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

        button.setTitleColor(.white, for: .normal)
        button.backgroundColor = PHColors.topaz

        button.layer.cornerRadius = 5
        button.layer.borderWidth = 1
        button.layer.borderColor = PHColors.topaz.cgColor

        button.setTitle("Continue", for: .normal)

        button.addTarget(
            self,
            action: #selector(continueButtonTapped),
            for: .touchUpInside)

        return button
    }()

    private let steps = TransactionFlowStep.allCases
    private var currentStepIndex = 0
    private let transaction: TransactionNotification

    // MARK: - Initialization

    init(transaction: TransactionNotification) {
        self.transaction = transaction

        super.init(nibName: nil, bundle: nil)
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        view.backgroundColor = PHColors.lightGray

        configureNavigationBar()
        configurePageViewController()
        configureBottomStackView()
        configureOptionsButton()

        // set initial viewcontroller
        continueFlow()
    }

    private func configureNavigationBar() {

        let stackView = UIStackView()
        stackView.translatesAutoresizingMaskIntoConstraints = false
        stackView.axis = .horizontal
        stackView.alignment = .center
        stackView.spacing = 5

        let image = UIImage(named: "phone_house")
        let imageView = UIImageView(image: image)
        imageView.contentMode = .center

        let label = UILabel()
        label.translatesAutoresizingMaskIntoConstraints = false
        label.text = "Phone House - Mobile subscription"
        label.textAlignment = .left
        label.font = PHFonts.wesBold(ofSize: 16)
        label.textColor = PHColors.greyishBrown

        stackView.widthAnchor.constraint(equalToConstant: UIScreen.main.bounds.width - 40).isActive = true

        stackView.addArrangedSubview(label)
        stackView.addArrangedSubview(imageView)

        navigationItem.titleView = stackView
        navigationController?.navigationBar.tintColor = PHColors.greyishBrown

        navigationItem.backBarButtonItem = UIBarButtonItem(
            title: "",
            style: .plain,
            target: nil,
            action: nil)

        NotificationCenter.default.addObserver(
            self,
            selector: #selector(pushWarningViewController),
            name: NSNotification.Name("show warning"),
            object: nil)

        NotificationCenter.default.addObserver(
            self,
            selector: #selector(pushVerificationViewController),
            name: NSNotification.Name("show verification"),
            object: nil)
    }

    private func configurePageViewController() {

        view.addSubview(pageViewController.view)
        pageViewController.view.translatesAutoresizingMaskIntoConstraints = false
        pageViewController.view.topAnchor.constraint(equalTo: view.safeAreaLayoutGuide.topAnchor).isActive = true
        pageViewController.view.leftAnchor.constraint(equalTo: view.safeAreaLayoutGuide.leftAnchor).isActive = true
        pageViewController.view.rightAnchor.constraint(equalTo: view.safeAreaLayoutGuide.rightAnchor).isActive = true
    }

    private func configureBottomStackView() {

        view.addSubview(bottomStackView)
        bottomStackView.heightAnchor.constraint(equalToConstant: 75).isActive = true
        bottomStackView.topAnchor.constraint(equalTo: pageViewController.view.bottomAnchor).isActive = true
        bottomStackView.rightAnchor.constraint(equalTo: pageViewController.view.rightAnchor, constant: -30).isActive = true
        bottomStackView.bottomAnchor.constraint(equalTo: view.safeAreaLayoutGuide.bottomAnchor).isActive = true

        bottomStackView.addArrangedSubview(declineButton)
        bottomStackView.addArrangedSubview(continueButton)
    }

    private func configureOptionsButton() {

        view.addSubview(optionsButton)
        optionsButton.leftAnchor.constraint(equalTo: view.leftAnchor, constant: 20).isActive = true
        optionsButton.centerYAnchor.constraint(equalTo: bottomStackView.centerYAnchor).isActive = true
    }

    // MARK: - PageViewController

    private func continueFlow() {

        guard currentStepIndex < steps.count else {
            respondToTransaction(isAccepted: true)
            return
        }

        // set initial viewcontroller
        let viewController = self.steps.map { $0.viewController(withTransaction: transaction) }[currentStepIndex]
        pageViewController.setViewControllers(
            [viewController],
            direction: .forward,
            animated: true)

        // increment for next step
        currentStepIndex += 1

        // last step!
        if currentStepIndex == steps.count {
            continueButton.setTitle("Finish!", for: .normal)
        }
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
                    default:
                        break
                    }
            }
        } catch {
            print(error)
        }
    }

    // MARK: - Selectors

    @objc private func pushWarningViewController(_ sender: Notification) {

        guard let viewController = sender.object as? UIViewController else {
            return
        }

        navigationController?.pushViewController(viewController, animated: true)
    }

    @objc private func pushVerificationViewController(_ sender: Notification) {

        guard let viewController = sender.object as? UIViewController else {
            return
        }

        navigationController?.pushViewController(viewController, animated: true)
    }

    @objc private func optionsButtonTapped(_ sender: UIButton) {

        let alert = UIAlertController(
            title: "External verification options",
            message: nil,
            preferredStyle: .actionSheet)

        alert.addAction(UIAlertAction(
            title: "Share with a contact",
            style: .default,
            handler: { _ in
                alert.dismiss(animated: true)
        }))

        alert.addAction(UIAlertAction(
            title: "Use audit tool",
            style: .default,
            handler: { _ in
                alert.dismiss(animated: true)
        }))

        alert.addAction(UIAlertAction(
            title: "Cancel",
            style: .cancel,
            handler: { _ in
                alert.dismiss(animated: true)
        }))

        present(alert, animated: true)
    }

    @objc private func declineButtonTapped(_ sender: UIButton) {
        presentCancellationAlert()
    }

    @objc private func continueButtonTapped(_ sender: UIButton) {
        continueFlow()
    }

    // MARK: - Helpers

    private func presentCancellationAlert() {

        let alert = UIAlertController(
            title: "Cancel transaction",
            message: "Are you sure you want to cancel?",
            preferredStyle: .alert)

        alert.addAction(UIAlertAction(
            title: "Yes",
            style: .destructive,
            handler: { [unowned self] _ in
                self.respondToTransaction(isAccepted: false)
        }))

        alert.addAction(UIAlertAction(
            title: "No",
            style: .default,
            handler: { _ in
                alert.dismiss(animated: true)
        }))

        present(alert, animated: true)
    }

    private func validateSelection() {

        let areAllItemsAccepted = transaction.items.filter { !$0.isAccepted }.count == 0
        continueButton.isEnabled = areAllItemsAccepted
    }
}
