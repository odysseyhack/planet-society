//
//  TransactionOverviewViewController.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

enum TransactionVerificationStatus {
    case verified, unverified

    var color: UIColor {
        switch self {
        case .verified:
            return PHColors.topaz
        case .unverified:
            return PHColors.red
        }
    }
}

enum PHTableViewViewCellType {
    case notification(
        type: TransactionNotificationType,
        text: String)
    case warning(text: String)
    case description(date: Date, title: String, description: String)
    case plugin(image: UIImage?, text: String)
    case transactionItem(item: TransactionItem)
    case selectionDisclosure(text: String)
    case selection(options: [String])
    case form(placeholder: String, text: String?)
}

protocol PHTableViewControllerDelegate: class {
    func didSelect(item: PHTableViewViewCellType)
    func didDeselect(item: PHTableViewViewCellType)
}

class PHTableViewController: UIViewController {

    // MARK: - Private properties

    lazy var tableView: UITableView = {

        let tableView = UITableView()
        tableView.translatesAutoresizingMaskIntoConstraints = false

        tableView.rowHeight = UITableView.automaticDimension
        tableView.estimatedRowHeight = 66
        tableView.separatorStyle = .none

        tableView.dataSource = self
        tableView.delegate = self

        tableView.register(
            UITableViewCell.self,
            forCellReuseIdentifier: String(describing: UITableViewCell.self))

        tableView.register(
            TransactionTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionTableViewCell.self))

        tableView.register(
            TransactionDescriptionTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionDescriptionTableViewCell.self))

        tableView.register(
            TransactionPluginTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionPluginTableViewCell.self))

        tableView.register(
            TransactionNotificationTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionNotificationTableViewCell.self))

        tableView.register(
            FormTextInputCell.self,
            forCellReuseIdentifier: String(describing: FormTextInputCell.self))

        return tableView
    }()

    private var items: [PHTableViewViewCellType]

    private let activityIndicatorViewController: PHActivityIndicatorViewController = {

        let viewController = PHActivityIndicatorViewController()
        viewController.modalPresentationStyle = .overCurrentContext
        viewController.modalTransitionStyle = .crossDissolve

        return viewController
    }()

    // MARK: - Properties

    weak var delegate: PHTableViewControllerDelegate?

    // MARK: - Initialization

    init(
        title: String? = nil,
        items: [PHTableViewViewCellType]) {

        self.items = items

        super.init(nibName: nil, bundle: nil)

        self.title = title
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        view.backgroundColor = PHColors.lightGray

        // configure navigation bar
        configureNavigationBar()
        configureTableView()
    }

    // MARK: - Configuration

    private func configureNavigationBar() {

        navigationItem.title = title
        navigationController?.navigationBar.tintColor = PHColors.greyishBrown
        navigationItem.backBarButtonItem = UIBarButtonItem(title: "", style: .plain, target: nil, action: nil)

        navigationController?.navigationBar.titleTextAttributes = [
            .font: PHFonts.wesBold(ofSize: 17),
            .foregroundColor: PHColors.greyishBrown
        ]
    }

    private func configureTableView() {

        view.addSubview(tableView)

        tableView.topAnchor.constraint(equalTo: view.topAnchor).isActive = true
        tableView.leftAnchor.constraint(equalTo: view.leftAnchor).isActive = true
        tableView.rightAnchor.constraint(equalTo: view.rightAnchor).isActive = true
        tableView.bottomAnchor.constraint(equalTo: view.bottomAnchor).isActive = true
    }
}

extension PHTableViewController: UITableViewDataSource {

    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return items.count
    }

    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {

        switch items[indexPath.row] {
        case .notification(let type, let text):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionNotificationTableViewCell.self),
                for: indexPath) as! TransactionNotificationTableViewCell

            cell.configure(withType: type, andText: text)

            return cell

        case .description(let date, let title, let description):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionDescriptionTableViewCell.self),
                for: indexPath) as! TransactionDescriptionTableViewCell

            cell.configure(withDate: date, andTitle: title, andDescription: description)

            return cell

        case .transactionItem(let item):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionTableViewCell.self),
                for: indexPath) as! TransactionTableViewCell

            let viewModel = TransactionTableViewCellViewModel(
                image: UIImage(),
                title: item.item,
                subtitle: (item.fields as NSArray).componentsJoined(by: ", "))
            cell.configure(withViewModel: viewModel)

            return cell

        case .warning(let text):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionTableViewCell.self),
                for: indexPath) as! TransactionTableViewCell

            let viewModel = TransactionTableViewCellViewModel(
                image: UIImage(),
                title: text,
                subtitle: "")
            cell.configure(withViewModel: viewModel)

            return cell

        case .plugin(let image, let text):
            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionPluginTableViewCell.self),
                for: indexPath) as! TransactionPluginTableViewCell

            cell.configure(
                withImage: image,
                andText: text,
                callback: { [unowned self] in
                    self.present(self.activityIndicatorViewController, animated: false)
                    self.perform(
                        #selector(self.dismissActivityIndicator),
                        with: nil,
                        afterDelay: 2)
            })

            return cell

        case .selectionDisclosure(let text):
            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: UITableViewCell.self),
                for: indexPath)

            cell.textLabel?.font = PHFonts.regular(ofSize: 14)
            cell.textLabel?.textColor = PHColors.greyishBrown
            cell.textLabel?.text = text
            cell.accessoryType = .disclosureIndicator

            return cell

        case .selection(let options):
            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: UITableViewCell.self),
                for: indexPath)

            return cell

        case .form(let placeholder, let text):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: FormTextInputCell.self),
                for: indexPath) as! FormTextInputCell

            cell.configure(
                withPlaceholder: placeholder,
                andText: text) { text in
                    print(text)
            }

            return cell
        }
    }

    // MARK: - Selectors

    @objc private func dismissActivityIndicator(_ sender: Any) {
        activityIndicatorViewController.dismiss(animated: true)

        items = [
            .notification(
                type: .verification,
                text: "This company is verified"),
            .description(
                date: Date(),
                title: "Personal details",
                description: "Please fill out your personal details."),
            .plugin(
                image: UIImage(named: "digid_button"),
                text: "Use the external DigiD plug-in to fill in your personal information (optional)."),
            .form(placeholder: "First name", text: "Gerard"),
            .form(placeholder: "Last name", text: "Huizinga"),
            .form(placeholder: "Date of birth", text: "04-11-1964"),
            .form(placeholder: "Address", text: "Weesperplein 43-2, Amsterdam"),
            .form(placeholder: "Email", text: "gerard.huizinga@gmail.com"),
            .form(placeholder: "BSN number", text: "264036232")
        ]

        tableView.reloadData()
    }
}

extension PHTableViewController: UITableViewDelegate {

    func tableView(_ tableView: UITableView, didSelectRowAt indexPath: IndexPath) {
        delegate?.didSelect(item: items[indexPath.row])

        switch items[indexPath.row] {
        case .transactionItem:
            tableView.cellForRow(at: indexPath)?.isSelected = true

        default:
            break
        }
    }
}
